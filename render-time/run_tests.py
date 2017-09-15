__author__ = 'creed'


def timed_test(command):
    import time
    start = time.time()
    os.system(command)
    return time.time() - start


def time_stamp():
    import datetime
    return datetime.datetime.now().strftime("%y%m%dT%H%M%S")


if __name__ == "__main__":
    import sys
    import os

    scene_name = sys.argv[1] if sys.argv[1:] else 'bmw27_cpu.blend'

    file_name = sys.argv[2] if sys.argv[2:] else 'resolutions.in'

    report_file_name = 'report_' + time_stamp() + '.out'
    report_file = open( report_file_name, 'w+' )

    with open(file_name) as f:
        lines = f.readlines()

    f = 0
    for line in lines:
        words = line.split()
        if len( words ) != 2:
            print 'Parse error: ' + line
        else:
            x = str( float( words[0] ) / 400. )
            y = str( float( words[1] ) / 300. )
            f += 1
            print 'x = ' + x + ', y = ' + y
            command  = 'blender '
            command += '-noaudio '
            command += '-b ' + scene_name + ' '
            command += '--python-expr '
#            command += '"import bpy; bpy.data.scenes[ 0 ].render.resolution_x = ' + x + '; bpy.data.scenes[ 0 ].render.resolution_y = ' + y + ';" '
            command += '"import bpy; bpy.data.scenes[ 0 ].render.resolution_x = 400; bpy.data.scenes[ 0 ].render.resolution_y = 300; ' \
                       'bpy.data.scenes[ 0 ].render.use_border = True; bpy.data.scenes[ 0 ].render.use_crop_to_border = True; ' \
                       'bpy.data.scenes[ 0 ].render.border_min_x = 0; bpy.data.scenes[ 0 ].render.border_min_y = 0; ' \
                       'bpy.data.scenes[ 0 ].render.border_max_x = ' + x + '; bpy.data.scenes[ 0 ].render.border_max_y = ' + y + ';" '
            command += '-f ' + str( f ) + ' '
            print 'command: ' + command
            time = timed_test( command )
            print 'python time: ' + str( time )
            report_file.write('{} {} {}\n'.format(words[0], words[1], time))
            report_file.flush()

    report_file.close()
